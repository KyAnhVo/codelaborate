#include "mainwindow.h"

#include <QLineEdit>
#include <QPlainTextEdit>
#include <QPushButton>
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QGridLayout>
#include <qboxlayout.h>

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
    
    // conn here and there
    connect(this->editor->document(), &QTextDocument::contentsChange,
            this->networkManager, &Network::sendUpdateToServer);
    connect(this->networkManager, &Network::receivedUpdate,
            this->editor, &Editor::update);

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

void MainWindow::createSession() {}
void MainWindow::joinSession() {}
void MainWindow::exitSession() {}
