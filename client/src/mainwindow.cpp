#include "mainwindow.h"
#include <qpushbutton.h>

MainWindow::MainWindow(QString serverIP, quint16 serverPort) {
    this->networkManager = new Network(serverIP, serverPort);

    this->createRoomButton  = new QPushButton("Create room");
    this->joinRoomButton    = new QPushButton("Join room");
    this->roomIDLabel       = new QLabel("Room ID");
}


