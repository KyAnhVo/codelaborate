#include "editor.h"
#include <qplaintextedit.h>

Editor::Editor(QWidget * parent) : QPlainTextEdit(parent) {}
Editor::~Editor() {}

/**
 * replaceStr = ''      => delete
 * charsReplaced = 0    => insert
 * Both                 => do nothing
 * Neither              => replace
 */
void Editor::applyOnlineChanges(int cursorPos, int charsReplaced, const QString& replaceStr) {
    QTextCursor editCursor(this->document());
    editCursor.setPosition(cursorPos);
    editCursor.setPosition(cursorPos + charsReplaced, QTextCursor::KeepAnchor);
    editCursor.insertText(replaceStr);
}
