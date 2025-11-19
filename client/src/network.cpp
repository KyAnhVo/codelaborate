#include "network.h"

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
    quint64 read_data = 0;
}
