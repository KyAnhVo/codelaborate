#include <QtWidgets>
#include <qapplication.h>

#include "mainwindow.h"

int main(int argc, char** argv) {
    QApplication app(argc, argv);

    MainWindow win;
    win.resize(1024, 768);
    win.setVisible(true);

    return app.exec();
}
