#include "editor.h"

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
        QTextCursor cursor(this->document());
        cursor.setPosition(position);
        cursor.movePosition(QTextCursor::Right, QTextCursor::KeepAnchor, insertLen);
        msg.insertStr = cursor.selectedText().toUtf8();
        msg.insertLen = msg.insertStr.length();
    } else {
        msg.insertLen = 0;
    }

    emit this->edited(msg);
}

void Editor::applyRemoteEdit(UpdateMsg msg) {
    qDebug() << "Remote change: " << msg.cursorPos
        << msg.deleteLen << msg.insertLen << msg.insertStr;
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
