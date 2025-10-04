#ifndef NETWORK_H
#define NETWORK_H

#include <qtmetamacros.h>
#include <string>
#include <QObject>

class Network : public QObject {
    Q_OBJECT

public:
    Network(int, std::string, int);
    ~Network();

public slots:
    void sendUpdateToServer(int, int, const QString&);

signals:
    void receivedUpdate(int, int, const QString&);


private:
    int         local_port;
    std::string server_ip;
    int         server_port;
};

#endif
