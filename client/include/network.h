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
    void entrySucceed();
    void entryFailed();
    void bogusSignal();
    
private:
    void recvEntryMsg(char);
    void recvUpdateMsg(char);

    QString     serverIP;
    quint16     serverPort;
    QTcpSocket  socket;
};

#endif
