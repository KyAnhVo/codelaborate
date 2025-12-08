// overall idea of the code:
// Let pendingOps be a queue of un-Acked msgs coming from our side to the server.
// Then any time we make a local change, we enqueue that op to pendingOps
// Then when the msg comes back, if it is our message then we know that
// the earliest un-Acked op is now acked, and we can dequeue it.
// Else if it is others' op we apply T(other, ours) for each of our op.
// Then we do the same: apply T(ours, other) on each of our op.

#include "editor.h"
#include "protocol.h"
#include <qtypes.h>

Editor::Editor() : QPlainTextEdit() {
    this->connect(
            this->document(), &QTextDocument::contentsChange,
            this, &Editor::onContentsChanged);
}

void Editor::onContentsChanged(int position, int deleteLen, int insertLen) {
    if (applyingRemoteEdit) return;

    qDebug() << "Editor change:" << position << deleteLen << insertLen;
    if (deleteLen == 0 && insertLen == 0) return;

    UpdateMsg msg;
    msg.op = MsgOp::UPDATE;
    msg.cursorPos = position;
    msg.deleteLen = deleteLen;

    if (insertLen > 0) {
        // get test from [position, position + insertLen)
        // Get raw text from document (preserves newlines)
        QString fullText = this->document()->toPlainText();
        QString text = fullText.mid(position, insertLen);
        
        msg.insertStr = text.toUtf8();
        msg.insertLen = msg.insertStr.size();
    } else {
        msg.insertLen = 0;
    }

    this->pendingOps.push_back(msg);
    
    emit this->edited(msg);
}

void Editor::applyRemoteEdit(UpdateMsg rawMsg, quint8 clientID) {
    if (clientID == this->clientID) {
        this->pendingOps.pop_front();
        return;
    }
    
    UpdateMsg msg = rawMsg;
    for (UpdateMsg& unackedMsg : this->pendingOps) {
        msg = this->transform(msg, unackedMsg, clientID, this->clientID);
        unackedMsg = this->transform(unackedMsg, rawMsg, this->clientID, clientID);
    }

    applyingRemoteEdit = true;
    this->blockSignals(true);
    this->document()->setUndoRedoEnabled(false);
    QTextCursor cursor(this->document());

    cursor.beginEditBlock();
    if (msg.deleteLen > 0) {
        cursor.setPosition(msg.cursorPos);
        cursor.movePosition(QTextCursor::Right, QTextCursor::KeepAnchor, msg.deleteLen);
        cursor.removeSelectedText();
    }

    if (msg.insertLen > 0) {
        cursor.setPosition(msg.cursorPos);
        cursor.insertText(QString::fromUtf8(msg.insertStr));
    }
    cursor.endEditBlock();

    this->document()->setUndoRedoEnabled(true);
    this->blockSignals(false);
    applyingRemoteEdit = false;
}

UpdateMsg Editor::transform(const UpdateMsg& target, const UpdateMsg& other, quint8 targetClientID, quint8 otherClientID) {
    UpdateMsg ret;
    ret.op = target.op;
    ret.insertLen = target.insertLen;
    ret.insertStr = target.insertStr;

    // first need to transform insertLen to insertLen in UTF-8 chars
    // for OT purpose
    quint64 targetInsertLenUtf8 = QString::fromUtf8(target.insertStr).length(),
            otherInsertLenUtf8  = QString::fromUtf8(other.insertStr).length();

    // then any cursorpos change for unintersected is based on these values
    quint64 targetAlterLen = targetInsertLenUtf8 - target.deleteLen,
            otherAlterLen = otherInsertLenUtf8 - other.deleteLen;
    
    // other op is earlier than target op in the document
    if (other.cursorPos < target.cursorPos) {
        // non-intersected previous op, intuitively we
        // push the cursor position by the amount of
        // alter the previous one was pushed up
        if (other.cursorPos + other.deleteLen < target.cursorPos) {
            ret.deleteLen = target.deleteLen;
            ret.cursorPos = target.cursorPos + otherAlterLen;
        }

        // intersected op: there exists 2 cases: either
        // inside case: [other start][target start][target end][other end], or
        // interweiving case: [other start][target start][other end][target end]
        else {
            // TODO: implement this
        }
    }

    // target op is earlier than other op in the document
    else if (target.cursorPos < other.cursorPos) {
        // non-intersected latter op, essentially the same
        if (target.cursorPos + target.deleteLen < other.cursorPos) {
            ret.deleteLen = target.deleteLen;
            ret.cursorPos = target.cursorPos;
        }

        // intersected op, there exists 2 cases: either
        // inside case: [target start][other start][other end][target end], or
        // interweiving case: [target start][other start][target end][other end]
        else {
            // TODO: implement this
        }
    }

    // target op at the same location as other op
    else {
        // If both are inserts, for a consistent cursor position,
        // we identify the clientID and whichever client has the lower clientID
        // we place it at the front.
        if (other.deleteLen == 0 && target.deleteLen == 0) {
            ret.deleteLen = 0;
            if (targetClientID < otherClientID) {
                ret.cursorPos = target.cursorPos;
            }
            else {
                ret.cursorPos = target.cursorPos + otherAlterLen;
            }
        }
    }

    return ret;
}
