#include "network.h"

#include <boost/asio.hpp>
#include <boost/asio/ip/tcp.hpp>
#include <boost/endian/conversion.hpp>
#include <iostream>

/**
 * UTIL FUNCTION DECLARATIONS
 */

template <typename T>
T stouint(std::string);

/**
 * CLASS FUNCTIONS
 */

Network::Network(std::string serverIpAddr, std::string serverPort) : io(), socket(io) {
    try {
        boost::asio::ip::tcp::resolver resolver(io);
        auto endpoints = resolver.resolve(serverIpAddr, serverPort);
        boost::asio::connect(this->socket, endpoints);
    }
    catch (std::exception& e) {
        std::cerr << "error: " << e.what() << std::endl;
        exit(404);
    }
}

Network::~Network() {}

void Network::connectToServer(Network::ConnType conn, std::string roomID) {
    uint32_t roomIDInt;
    roomIDInt = stouint<uint32_t>(roomID);
    roomIDInt = boost::endian::native_to_big(roomIDInt);

    uint8_t buffer[5];
    switch (conn) {
        case Network::ConnType::CREATE:
            buffer[0] = 'C';
            break;
        case Network::ConnType::JOIN:
            buffer[0] = 'J';
            break;
    }
    memcpy(buffer + 1, &roomIDInt, sizeof(uint32_t));
    boost::asio::write(this->socket, boost::asio::buffer(buffer, 5));
}

void Network::sendUpdateToServer(int pos, int deleteLen, const QString& addStr) {
    
}

/**
 * UTIL FUNCTION DEFINITIONS
 */

template <typename T>
T stouint(std::string str) {
    unsigned long temp = std::stoul(str);
    return static_cast<T>(temp);
}
