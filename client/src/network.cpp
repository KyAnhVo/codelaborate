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
    buf.append(msg.insertStr.toUtf8());

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
    switch (buf[0]) {
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
}

void Network::recvEntryMsg(char msgStatus) {

}
