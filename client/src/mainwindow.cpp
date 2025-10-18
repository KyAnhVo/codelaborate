#include "mainwindow.h"

#include <QLineEdit>
#include <QPlainTextEdit>
#include <QPushButton>
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QGridLayout>
#include <QBoxLayout>
#include <qnamespace.h>

#define SERVER_ADDR_TEST "127.0.0.1"
#define SERVER_ADDR_PROD "127.0.0.1"    // subject to change
#define SERVER_PORT_TEST 80
#define SERVER_PORT_PROD 80             // subject to change

MainWindow::MainWindow(QWidget * parent) : QMainWindow(parent) {
    // buttons setup
    this->createSessionButton   = new QPushButton("Create Session");
    this->joinSessionButton     = new QPushButton("Join Session");
    this->exitSessionButton     = new QPushButton("Exit current session");
    
    // editor conf
    this->editor = new Editor();

    // session id link
    this->sessionIdLineEdit = new QLineEdit();

    // network manager
    this->networkManager = new Network("127.0.0.1", "80");
    this->networkManager->moveToThread(&this->networkThread);
    
    // Thread startup and cleanup
    connect(&this->networkThread, &QThread::finished,
            this->networkManager, &Network::deleteLater);

    // Converts from (pos, delLen, insertLen) to (pos, delLen, insertStr)
    connect(this->editor->document(), &QTextDocument::contentsChange,
            this, &MainWindow::receiveUpdateLens,
            Qt::QueuedConnection);

    // sends the updated (pos, delLen, insertStr) to server
    connect(this, &MainWindow::sendReplacementInfo,
            this->networkManager, &Network::sendUpdateToServer,
            Qt::QueuedConnection);

    // receives updates from server in (pos, delLen, insertStr)
    // then updates the editor with the corresponding update
    connect(this->networkManager, &Network::receivedUpdate,
            this->editor, &Editor::applyOnlineChanges,
            Qt::QueuedConnection);

    // general sends create/join to networkManager
    connect(this, &MainWindow::connectToServer,
            this->networkManager, &Network::connectToServer);

    // Sends create session to server
    connect(this->createSessionButton, &QPushButton::clicked,
            this, &MainWindow::createSessionClicked);

    // Sends join session to server
    connect(this->joinSessionButton, &QPushButton::clicked,
            this, &MainWindow::joinSessionClicked);
    

    // Session Box
    QHBoxLayout * sessionBox = new QHBoxLayout;
    sessionBox->addWidget(this->createSessionButton);
    sessionBox->addWidget(this->sessionIdLineEdit);
    sessionBox->addWidget(this->joinSessionButton);
    sessionBox->addWidget(this->exitSessionButton);

    // Main layout
    QVBoxLayout * mainLayout = new QVBoxLayout;
    mainLayout->addLayout(sessionBox);
    mainLayout->addWidget(this->editor);
    QWidget * centralWidget = new QWidget();
    centralWidget->setLayout(mainLayout);
    this->setCentralWidget(centralWidget);
}

MainWindow::~MainWindow() {}

void MainWindow::receiveUpdateLens(int pos, int deleteLen, int insertLen) {
    QTextCursor cursor(this->editor->document());
    cursor.setPosition(pos);
    cursor.setPosition(pos + insertLen - 1, QTextCursor::KeepAnchor);
    const QString addedStr = cursor.selectedText();
    emit sendReplacementInfo(pos, deleteLen, addedStr);
}

void MainWindow::createSessionClicked() {
    QString id = this->sessionIdLineEdit->text();
    std::string idStd = id.toStdString();
    emit this->connectToServer(Network::ConnType::CREATE, idStd);
}

void MainWindow::joinSessionClicked() {
    QString id = this->sessionIdLineEdit->text();
    std::string idStd = id.toStdString();
    emit this->connectToServer(Network::ConnType::JOIN, idStd);
}
void MainWindow::exitSessionClicked() {}
