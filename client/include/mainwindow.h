#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QPushButton>
#include <QLineEdit>
#include <QPlainTextEdit>
#include <QThread>
#include <qtmetamacros.h>

#include "editor.h"
#include "network.h"

class MainWindow : public QMainWindow {
    Q_OBJECT

public:
    MainWindow(QWidget * parent=nullptr);
    ~MainWindow();

signals:
    void sendReplacementInfo(int pos, int deleteLen, const QString& insertStr);

public slots:
    // buttons clicked slots
    void createSession();
    void joinSession();
    void exitSession();

    // editor update slots
    void receiveUpdateLens(int, int, int);

private:
    bool has_session;
    QThread networkThread;

    QPushButton * createSessionButton,
                * joinSessionButton,
                * exitSessionButton;

    QLineEdit * sessionIdLineEdit;
    
    Editor * editor;

    Network * networkManager;
};

#endif
