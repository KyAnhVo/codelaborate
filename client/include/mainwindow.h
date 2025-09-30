#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QPushButton>
#include <QLineEdit>
#include <QPlainTextEdit>
#include <qtmetamacros.h>

#include "editor.h"

class MainWindow : public QMainWindow {
    Q_OBJECT

public:
    MainWindow(QWidget * parent=nullptr);
    ~MainWindow();

public slots:
    // buttons clicked slots
    void createSession();
    void joinSession();
    void exitSession();

private:
    bool has_session;

    QPushButton * createSessionButton,
                * joinSessionButton,
                * exitSessionButton;

    QLineEdit * sessionId;
    
    Editor * editor;
};

#endif
