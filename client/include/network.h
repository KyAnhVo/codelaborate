#ifndef NETWORK_H
#define NETWORK_H

#include <QTcpSocket>
#include <QObject>
#include <QtTypes>

#include "protocol.h"

class Network : public QObject {
    Q_OBJECT
public:
    explicit Network(QString serverIP, quint16 serverPort = 80);

public slots:
    void sendUpdateMsg(UpdateMsg);
    void sendEntryMsg(EntryMsg);
    void recvMsg();

signals:
    void updateMsgArrived(UpdateMsg);
    void closeConnMsgArrived();
    void entrySucceed(quint32 roomID);
    void entryFailed();
    void bogusSignal();
    
private:
    // helper functions to receive message

    void recvEntryMsg(char);
    void recvUpdateMsg(char);

    // data extraction from socket

    template <typename t>
    t recvUnsignedIntOfType();
    QByteArray readStr(quint64);

    QString     serverIP;
    quint16     serverPort;
    QTcpSocket  socket;
};

#endif
