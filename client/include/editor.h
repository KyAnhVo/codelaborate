#ifndef EDITOR_H
#define EDITOR_H

#include <QPlainTextEdit>
#include <QtWidgets>
#include <deque>

#include "protocol.h"

class Editor : public QPlainTextEdit {
    Q_OBJECT
public:
    explicit Editor();
    quint8 clientID;

public slots:
    void onContentsChanged(int position, int deleteLen, int insertLen);
    void applyRemoteEdit(UpdateMsg msg, quint8 clientID);

signals:
    void edited(UpdateMsg msg);
private:
    bool applyingRemoteEdit = false;
    std::deque<UpdateMsg> pendingOps;
    UpdateMsg transform(const UpdateMsg& targetOp, const UpdateMsg& otherOp);
};

#endif
