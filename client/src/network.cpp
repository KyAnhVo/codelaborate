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

void Network::sendUpdateMsg(MsgOp op, quint64 cursorPos, quint64 deleteLen, quint64 insertLen, QString& insertStr) {
    QByteArray buf;

    // append MsgOp
    switch(op) {
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
    cursorPos = qToBigEndian(cursorPos);
    buf.append(reinterpret_cast<char*>(&cursorPos), sizeof(cursorPos));

    // append deleteLen
    deleteLen = qToBigEndian(deleteLen);
    buf.append(reinterpret_cast<char*>(&deleteLen), sizeof(deleteLen));

    // append insertLen
    insertLen = qToBigEndian(insertLen);
    buf.append(reinterpret_cast<char*>(&insertLen), sizeof(insertLen));

    // append insertStr
    buf.append(insertStr.toUtf8());

    this->socket.write(buf);
}

void Network::sendEntryMsg(MsgOp op, quint32 roomID) {
    QByteArray buf;

    // append op
    switch(op) {
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
    roomID = qToBigEndian(roomID);
    buf.append(reinterpret_cast<char*>(&roomID), sizeof(roomID));

    this->socket.write(buf);
}
