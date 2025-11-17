#ifndef NETWORK_H
#define NETWORK_H

#include <QTcpSocket>
#include <qobject.h>
#include <qtcpsocket.h>
#include <qtmetamacros.h>

enum class MsgOp {
    CREATE,
    JOIN,
    CLOSE_CONN,
    UPDATE,
};

class Network : public QObject {
    Q_OBJECT
public:
    explicit Network(QString serverIP, quint16 serverPort = 80);

public slots:
    void sendUpdateMsg(MsgOp op, quint64 cursorPos, quint64 deleteLen, quint64 insertLen, QString& insertStr);
    void sendEntryMsg(MsgOp op, quint32 roomID);
    
private:
    QString     serverIP;
    quint16     serverPort;
    QTcpSocket  socket;
};

#endif
