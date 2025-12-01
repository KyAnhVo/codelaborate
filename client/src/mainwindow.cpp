#include "mainwindow.h"
#include "editor.h"

#include <iostream>
#include <QtWidgets>
#include <qboxlayout.h>
#include <qwidget.h>

MainWindow::MainWindow(QString serverIP, quint16 serverPort) {
    // Initiate objects
    this->networkManager    = new Network(serverIP, serverPort);
    this->createRoomButton  = new QPushButton("Create room");
    this->joinRoomButton    = new QPushButton("Join room");
    this->exitRoomButton    = new QPushButton("Exit room");
    this->roomIDLabel       = new QLabel("Room ID");
    this->roomIDLineEdit    = new QLineEdit();
    this->editor            = new Editor();

    // metadata
    this->connected = false;

    // widget inits
    this->editor->setReadOnly(true);
    this->exitRoomButton->setDisabled(true);

    this->connectObjs();
    this->setupGui();
}

void MainWindow::connectObjs() {
    // connect objects: entry room
    connect(this->createRoomButton, &QPushButton::clicked,
            this, &MainWindow::createRoom);
    connect(this->joinRoomButton, &QPushButton::clicked,
            this, &MainWindow::joinRoom);
    connect(this, &MainWindow::sendEntryMsg,
            this->networkManager, &Network::sendEntryMsg);
    connect(this->networkManager, &Network::entrySucceed,
            this, &MainWindow::onJoinRoomSucceed);
    connect(this->networkManager, &Network::entryFailed,
            this, &MainWindow::onJoinRoomFailed);

    // connect objects: disconnection
    connect(this->networkManager, &Network::closeConnMsgArrived,
            this, &MainWindow::disconnect);

    // connect objects: editor
    connect(this->editor, &Editor::edited,
            this->networkManager, &Network::sendUpdateMsg);
    connect(this->networkManager, &Network::updateMsgArrived,
            this->editor, &Editor::applyRemoteEdit);
}


void MainWindow::setupGui() {
    /** SETUP:
     * RoomID <line edit> <createRoomButton> <joinRoomButton> <exitRoomButton>
     * The actual editor
     */

    QHBoxLayout * metadataBox = new QHBoxLayout();
    metadataBox->addWidget(this->roomIDLabel, 0);
    metadataBox->addWidget(this->roomIDLineEdit, 1);
    metadataBox->addWidget(this->createRoomButton, 0);
    metadataBox->addWidget(this->joinRoomButton, 0);
    metadataBox->addWidget(this->exitRoomButton, 0);

    QVBoxLayout * overallBox = new QVBoxLayout();
    overallBox->addLayout(metadataBox, 0);
    overallBox->addWidget(this->editor, 1);

    QWidget * central = new QWidget;
    central->setLayout(overallBox);
    this->setCentralWidget(central);
}

void MainWindow::joinRoom() {
    std::cout << "enter joinRoom()" << std::endl;
    QString roomID = this->roomIDLineEdit->text();
    EntryMsg msg;
    msg.roomID = static_cast<quint32>(roomID.toUInt());
    msg.op = MsgOp::JOIN;
    std::cout << "Joining room, roomID = " << msg.roomID << std::endl;
    emit this->sendEntryMsg(msg);
    std::cout << "exit joinRoom()" << std::endl;
}

void MainWindow::createRoom() {
    std::cout << "enter createRoom()" << std::endl;
    EntryMsg msg;
    msg.roomID = 0; // room id does not matter
    msg.op = MsgOp::CREATE;
    emit this->sendEntryMsg(msg);
    std::cout << "exit createRoom()" << std::endl;
}

void MainWindow::onJoinRoomSucceed(quint32 roomID, quint8 clientID) {
    this->roomIDLineEdit->setText(QString::number(roomID));
    this->roomIDLineEdit->setReadOnly(true);
    this->exitRoomButton->setDisabled(false);
    this->joinRoomButton->setDisabled(true);
    this->createRoomButton->setDisabled(true);
    this->editor->setReadOnly(false);
}

void MainWindow::onJoinRoomFailed() {

}

void MainWindow::disconnect() {
    this->roomIDLineEdit->setText("");
    this->roomIDLineEdit->setReadOnly(false);
    this->exitRoomButton->setDisabled(true);
    this->joinRoomButton->setDisabled(false);
    this->createRoomButton->setDisabled(false);
    this->editor->clear();
    this->editor->setReadOnly(true);
}
