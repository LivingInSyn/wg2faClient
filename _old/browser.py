import sys
from PyQt5.QtCore import QUrl
from PyQt5.QtWidgets import QWidget, QApplication, QVBoxLayout
from PyQt5.QtWidgets import QApplication
from PyQt5.QtWebEngineWidgets import QWebEngineView

class Browser(QWidget):
    # based on code from: https://zetcode.com/pyqt/qwebengineview/
    def __init__(self, url):
        super().__init__()
        self.url = url
        self.initUI()

    def initUI(self):
        vbox = QVBoxLayout(self)
        self.webEngineView = QWebEngineView()
        self.loadPage()
        vbox.addWidget(self.webEngineView)
        self.setLayout(vbox)
        self.setGeometry(300, 300, 350, 250)
        self.setWindowTitle('QWebEngineView')
        self.show()

    def loadPage(self):
        url = QUrl.fromUserInput(self.url)
        self.webEngineView.load(url)
