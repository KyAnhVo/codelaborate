#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QPushButton>
#include <QLineEdit>
#include <qtmetamacros.h>

class MainWindow : public QMainWindow {
    Q_OBJECT

public:
    MainWindow(QWidget * parent=nullptr);
    ~MainWindow();

public slots:
    // buttons clicked slots
    void on_create_session_button_clicked();
    void on_join_session_button_clicked();
    void on_exit_session_button_clicked();

private:
    bool has_session;

    QPushButton * create_session_button,
                * join_session_button,
                * exit_session_button;

    QLineEdit * session_id;
};

#endif
