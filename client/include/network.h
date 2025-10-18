#ifndef NETWORK_H
#define NETWORK_H

#include <boost/asio/io_context.hpp>
#include <qtmetamacros.h>
#include <string>
#include <QObject>
#include <boost/asio.hpp>

class Network : public QObject {
    Q_OBJECT

public:
    // used to specify connection type
    enum class ConnType { CREATE, JOIN, };

    Network(std::string, std::string);
    ~Network();

public slots:
    void sendUpdateToServer(int, int, const QString&);
    void connectToServer(Network::ConnType conn, std::string roomID);

signals:
    void receivedUpdate(int, int, const QString&);


private:
    boost::asio::io_context io;
    boost::asio::ip::tcp::socket socket;
};

#endif
