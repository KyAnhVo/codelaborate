#include "network.h"
#include "protocol.h"

#include <QtEndian>
#include <qendian.h>
#include <qobject.h>
#include <stdexcept>

Network::Network(QString serverIP, quint16 serverPort) {
    this->serverPort = serverPort;
    this->serverIP = serverIP;
    this->socket.connectToHost(serverIP, serverPort);

    connect(&(this->socket), &QTcpSocket::readyRead,
            this, &Network::recvMsg);
}

void Network::sendUpdateMsg(UpdateMsg msg) {
    QByteArray buf;

    // append MsgOp
    switch(msg.op) {
        case MsgOp::CLOSE_CONN:
            buf.append(static_cast<quint8>(0));
            break;
        case MsgOp::UPDATE:
            buf.append(static_cast<quint8>(1));
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
    quint64 readData = 0;
    UpdateMsg msg;
    char buf[8];

    readData = this->socket.read(buf, 1);
    switch (buf[0]) { // not convert this to detect faulty msg status (4 to 7 might appear?)
        case static_cast<char>(MsgStatus::ENTRY_ERR):
        case static_cast<char>(MsgStatus::ENTRY_OK):
            this->recvEntryMsg(buf[0]);
            break;
        case static_cast<char>(MsgStatus::CLOSE_CONN):
        case static_cast<char>(MsgStatus::UPDATE):
            this->recvUpdateMsg(buf[0]);
            break;
        default:
            emit this->bogusSignal();
    }
}

void Network::recvUpdateMsg(char msgStatus) {
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
    msg.deleteLen = this->recvUnsignedIntOfType<quint64>();
    msg.insertLen = this->recvUnsignedIntOfType<quint64>();
    msg.insertStr = this->readStr(msg.insertLen);

    emit this->updateMsgArrived(msg);
}

void Network::recvEntryMsg(char msgStatus) {
    switch (static_cast<MsgStatus>(msgStatus)) {
        case MsgStatus::ENTRY_OK:
            emit this->entrySucceed();
            break;
        case MsgStatus::ENTRY_ERR:
            emit this->entryFailed();
            break;
        default: // never happens
            throw std::runtime_error("Wrong function");
    }
}

template <typename t>
t Network::recvUnsignedIntOfType() {
    static_assert(std::is_unsigned<t>::value && std::is_integral<t>::value,
              "recvUnsignedIntOfType requires an unsigned integer type");
    t val = 0;
    qint64 dataSize = sizeof(t);
    qint64 dataRead = 0;
    while (dataRead < dataSize) {
        dataRead += this->socket.read(
                reinterpret_cast<char*>(&val) + dataRead, dataSize - dataRead);
    }
    return qFromBigEndian(val);
}

QByteArray Network::readStr(quint64 byteCount) {
    quint64 bytesRead = 0;
    char* buf = new char[byteCount];
    while (bytesRead < byteCount) {
        qint64 currRead = this->socket.read(
                buf + bytesRead, byteCount - bytesRead);
        bytesRead += currRead;
    }
    QByteArray str = QByteArray(buf, byteCount);
    return str;
}
