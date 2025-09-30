#ifndef EDITOR_H
#define EDITOR_H

#include <QPlainTextEdit>
#include <qtmetamacros.h>

class Editor : public QPlainTextEdit {
    Q_OBJECT

public:
    Editor(QWidget * parent = nullptr);
    ~Editor();

public slots:
    void update(int, int, const QString&);
};

#endif
