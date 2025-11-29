#ifndef EDITOR_H
#define EDITOR_H

#include <QPlainTextEdit>
#include "protocol.h"

class Editor : public QPlainTextEdit {
    Q_OBJECT
public:
    explicit Editor();

public slots:
    void onContentsChanged(int position, int deleteLen, int insertLen);
    void applyRemoteEdit(UpdateMsg msg);

signals:
    void edited(UpdateMsg msg);
private:
    bool applyingRemoteEdit = false;
};

#endif
