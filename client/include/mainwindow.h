#ifndef MAIN_WINDOW_H
#define MAIN_WINDOW_H

#include "network.h"
#include "protocol.h"
#include "editor.h"

#include <QMainWindow>
#include <QPushButton>
#include <QPlainTextEdit>
#include <QLabel>
#include <QLineEdit>

class MainWindow : public QMainWindow {
    Q_OBJECT

public:
    explicit MainWindow(QString, quint16);

public slots:
    void joinRoom();
    void createRoom();
    void disconnect();
    void onJoinRoomSucceed(quint32);
    void onJoinRoomFailed();

signals:
    void sendEntryMsg(EntryMsg msg);

private:
    // Qt GUI objects
    QPushButton     * joinRoomButton,
                    * createRoomButton,
                    * exitRoomButton;
    QLabel          * roomIDLabel;
    QLineEdit       * roomIDLineEdit;
    Editor          * editor;
    Network         * networkManager;

    // status values
    bool            connected;

    // helper functions
    void connectObjs();
    void setupGui();
};

#endif
