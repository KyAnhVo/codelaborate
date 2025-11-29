#ifndef MAIN_WINDOW_H
#define MAIN_WINDOW_H

#include "network.h"
#include "protocol.h"

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

signals:


private:
    QPushButton     * joinRoomButton,
                    * createRoomButton;
    QLabel          * roomIDLabel;
    QLineEdit       * roomIDLineEdit;
    QPlainTextEdit  * editor;
    Network         * networkManager;
};

#endif
