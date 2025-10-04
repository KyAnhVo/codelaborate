#include "mainwindow.h"

#include <QLineEdit>
#include <QPlainTextEdit>
#include <QPushButton>
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QGridLayout>
#include <QBoxLayout>
#include <qnamespace.h>

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
    this->networkManager = new Network(1900, "loopback", 80);
    this->networkManager->moveToThread(&this->networkThread);
    
    // Thread startup and cleanup
    connect(&this->networkThread, &QThread::finished,
            this->networkManager, &Network::deleteLater);

    // Param for current slot/signal: int, int, int
    // this one essentially converts from (pos, delLen, insertLen)
    // to (pos, delLen, insertStr)
    connect(this->editor->document(), &QTextDocument::contentsChange,
            this, &MainWindow::receiveUpdateLens,
            Qt::QueuedConnection);
    // Param for current slot/signal: int, int, const QString&
    // this one essentially sends the updated (pos, delLen, insertStr)
    // to server
    connect(this, &MainWindow::sendReplacementInfo,
            this->networkManager, &Network::sendUpdateToServer,
            Qt::QueuedConnection);
    // Param for current slot/signal: int, int, const QString&
    // this one receives updates from server in (pos, delLen, insertStr)
    // then updates the editor with the corresponding update
    connect(this->networkManager, &Network::receivedUpdate,
            this->editor, &Editor::applyOnlineChanges,
            Qt::QueuedConnection);

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
void MainWindow::createSession() {}
void MainWindow::joinSession() {}
void MainWindow::exitSession() {}
