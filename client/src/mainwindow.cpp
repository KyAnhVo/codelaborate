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
}
