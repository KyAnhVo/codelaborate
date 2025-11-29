#include "network.h"
#include "protocol.h"

#include <QtEndian>
#include <qendian.h>
#include <qlogging.h>
#include <qobject.h>
#include <QDebug>
#include <stdexcept>

Network::Network(QString serverIP, quint16 serverPort) {
    this->serverPort = serverPort;
    this->serverIP = serverIP;
    this->socket.connectToHost(serverIP, serverPort);

    // Wait for connection (blocking, but simple for MVP)
    if (!this->socket.waitForConnected(3000)) {  // 3 second timeout
        qWarning() << "Failed to connect:" << this->socket.errorString();
    }
    this->socket.setSocketOption(QAbstractSocket::LowDelayOption, 1);
    this->socket.setSocketOption(QAbstractSocket::KeepAliveOption, 1);

    connect(&(this->socket), &QTcpSocket::readyRead,
            this, &Network::recvMsg);
}

void Network::sendUpdateMsg(UpdateMsg msg) {
    QByteArray buf;

    // append MsgOp
    switch(msg.op) {
        case MsgOp::CLOSE_CONN:
            buf.append(static_cast<quint8>(1));
            break;
        case MsgOp::UPDATE:
            buf.append(static_cast<quint8>(0));
            break;
        default:
            throw std::runtime_error("Operation invalid for Update msg");
    }

    // append cursorPos
    quint64 cursorPos = qToBigEndian(msg.cursorPos);
    buf.append(reinterpret_cast<char*>(&cursorPos), sizeof(cursorPos));

    // append deleteLen
    quint64 deleteLen = qToBigEndian(msg.deleteLen);
    buf.append(reinterpret_cast<char*>(&deleteLen), sizeof(deleteLen));

    // append insertLen
    quint64 insertLen = qToBigEndian(msg.insertLen);
    buf.append(reinterpret_cast<char*>(&insertLen), sizeof(insertLen));

    // append insertStr
    buf.append(msg.insertStr);

    this->socket.write(buf);
}

void Network::sendEntryMsg(EntryMsg msg) {
    QByteArray buf;

    // append op
    switch(msg.op) {
        case MsgOp::JOIN:
            buf.append('J');
            break;
        case MsgOp::CREATE:
            buf.append('C');
            break;
        default:
            throw std::runtime_error("Operation invalid for Entry message");
    }

    // append roomID
    quint32 roomID = qToBigEndian(msg.roomID);
    buf.append(reinterpret_cast<char*>(&roomID), sizeof(roomID));

    this->socket.write(buf);
}

void Network::recvMsg() {
    qDebug() << "Enter recvMsg()";
    quint64 readData = 0;
    UpdateMsg msg;
    char buf;

    readData = this->socket.read(&buf, 1);
    switch (buf) { // not convert this to detect faulty msg status (4 to 7 might appear?)
        case static_cast<char>(MsgStatus::ENTRY_ERR):
        case static_cast<char>(MsgStatus::ENTRY_OK):
            this->recvEntryMsg(buf);
            break;
        case static_cast<char>(MsgStatus::CLOSE_CONN):
        case static_cast<char>(MsgStatus::UPDATE):
            this->recvUpdateMsg(buf);
            break;
        default:
            emit this->bogusSignal();
    }
    qDebug() << "Exit recvMsg()";
}

void Network::recvUpdateMsg(char msgStatus) {
    qDebug() << "Enter recvUpdateMsg()";
    UpdateMsg msg;

    switch (static_cast<MsgStatus>(msgStatus)) {
        case MsgStatus::CLOSE_CONN:
            msg.op = MsgOp::CLOSE_CONN;
            break;
        case MsgStatus::UPDATE:
            msg.op = MsgOp::UPDATE;
            break;
        default: // never happens
            throw std::runtime_error("Wrong function");
    }

    msg.cursorPos = this->recvUnsignedIntOfType<quint64>();
    qDebug() << "CursorPos: " << msg.cursorPos;
    msg.deleteLen = this->recvUnsignedIntOfType<quint64>();
    qDebug() << "deleteLen: " << msg.deleteLen;
    msg.insertLen = this->recvUnsignedIntOfType<quint64>();
    qDebug() << "insertLen: " << msg.insertLen;
    msg.insertStr = this->readStr(msg.insertLen);
    qDebug() << "insertStr: " << QString::fromUtf8(msg.insertStr);

    if (msg.op == MsgOp::CLOSE_CONN)
        emit this->closeConnMsgArrived();
    else
        emit this->updateMsgArrived(msg);
    qDebug() << "Exit recvUpdateMsg()";
}

void Network::recvEntryMsg(char msgStatus) {
    qDebug() << "Enter recvEntryMsg()";
    switch (static_cast<MsgStatus>(msgStatus)) {
        case MsgStatus::ENTRY_OK: {
                qDebug() << "Entry successful";
                quint32 roomID = this->recvUnsignedIntOfType<quint32>();
                qDebug() << "Entry successful for roomID "
                    << roomID;
                emit this->entrySucceed(roomID);
                break;
            }
        case MsgStatus::ENTRY_ERR:
            qDebug() << "Entry failed";
            emit this->entryFailed();
            break;
        default: // never happens
            throw std::runtime_error("Wrong function");
    }
    qDebug() << "Exit recvEntryMsg()";
}

template <typename t>
t Network::recvUnsignedIntOfType() {
    static_assert(std::is_unsigned<t>::value && std::is_integral<t>::value,
              "recvUnsignedIntOfType requires an unsigned integer type");
    t val = 0;
    qint64 dataSize = sizeof(t);
    qint64 dataRead = 0;
    while (dataRead < dataSize) {
        // Wait for data
        if (this->socket.bytesAvailable() < (dataSize - dataRead)) {
            if (!this->socket.waitForReadyRead(3000)) {
                qWarning() << "Timeout reading integer";
                return 0;
            }
        }
        
        qint64 justRead = this->socket.read(
                reinterpret_cast<char*>(&val) + dataRead, dataSize - dataRead);
        dataRead += justRead;
    }
    return qFromBigEndian(val);
}

QByteArray Network::readStr(quint64 byteCount) {
    quint64 bytesRead = 0;
    char* buf = new char[byteCount];
    while (bytesRead < byteCount) {
        // Wait for data to be available
        if (this->socket.bytesAvailable() == 0) {
            if (!this->socket.waitForReadyRead(3000)) {  // 3 second timeout
                qWarning() << "Timeout waiting for" << byteCount << "bytes, got" << bytesRead;
                delete[] buf;
                return QByteArray();
            }
        }
        qint64 currRead = this->socket.read(
                buf + bytesRead, byteCount - bytesRead);
        bytesRead += currRead;
    }
    QByteArray str = QByteArray(buf, byteCount);
    return str;
}
