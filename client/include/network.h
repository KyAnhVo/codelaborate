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
private:
    int         local_port;
    std::string server_ip;
    int         server_port;
};

#endif
