#include <QtWidgets>
#include <qapplication.h>

#include <iostream>

#include "mainwindow.h"

int main(int argc, char** argv) {
    QApplication app(argc, argv);

    if (argc != 3) {
        std::cout << "Usage: codelaborate <serverIP> <serverPort>" << std::endl;
        return 1;
    }

    QString serverIP(argv[1]);
    quint16 serverPort = static_cast<uint16_t>(std::stoi(argv[2]));

    MainWindow win(serverIP, serverPort);

    win.resize(1024, 768);
    win.setVisible(true);
    return app.exec();
}
