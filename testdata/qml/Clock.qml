import QtQuick
import Quickshell

Scope {
  id: root

  property string time: Qt.formatTime(new Date(), "hh:mm")

  Timer {
    interval: 1000
    running: true
    repeat: true
    onTriggered: root.time = Qt.formatTime(new Date(), "hh:mm")
  }
}
