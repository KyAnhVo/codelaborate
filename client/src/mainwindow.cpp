#include "mainwindow.h"
#include <qline.h>
#include <qplaintextedit.h>
#include <qpushbutton.h>

MainWindow::MainWindow(QWidget * parent) : QMainWindow(parent) {
    // buttons setup
    this->createSessionButton   = new QPushButton("Create Session");
    this->joinSessionButton     = new QPushButton("Join Session");
    this->exitSessionButton     = new QPushButton("Exit current session");
    
    // editor conf
    this->editor = new Editor();

    // session id link
    this->sessionId = new QLineEdit();
    
    // conn here and there
    connect(this->editor->document(), &QTextDocument::contentsChange,
            this->networkManager, &Network::sendUpdateToServer);
    connect(this->networkManager, &Network::receivedUpdate,
            this->editor, &Editor::update);
}
