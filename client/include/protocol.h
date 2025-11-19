#ifndef PROTOCOL_H
#define PROTOCOL_H

#include <QtTypes>
#include <QString>

enum class MsgOp : uint8_t {
    CREATE,
    JOIN,
    CLOSE_CONN,
    UPDATE,
};

enum class EntryStatus : uint8_t {
    OK,
    ERROR,
};

struct UpdateMsg {
    MsgOp op;
    quint64 cursorPos;
    quint64 deleteLen;
    quint64 insertLen;
    QString insertStr;
};

struct EntryMsg {
    MsgOp   op;
    quint32 roomID;
};

#endif
