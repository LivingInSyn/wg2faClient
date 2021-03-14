import secrets
import base64
from hashlib import sha256
# for http server
import http.cookies
from http.server import BaseHTTPRequestHandler, HTTPServer
from urllib.parse import urlparse, parse_qs, quote_plus
import time
# for the GUI
import webbrowser
import json
#
from wg2fa_client import Wg2faClient

BASE_DOMAIN='dev-318981.oktapreview.com'
DOMAIN=f'{BASE_DOMAIN}/oauth2/default/v1'
REDIRECT_URL='http://localhost:8080/implicit/callback'
CLIENT_ID='0oawepftsdT43o2CM0h7'

CALLBACK_ARGS=None
STATE = ''
VERIFIER = ''
TOKEN_CALLBACK = check_token

class HttpCallbackHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        self.send_response(200)
        self.send_header("Content-type", "text/html")
        self.end_headers()
        self.data_string = self.rfile.read(int(self.headers['Content-Length']))
        data = json.loads(self.data_string)
        if 'access_token' in data:
            TOKEN_CALLBACK(data)
        print(data)
        # TODO: post to wg2fa

    def do_DELETE(self):
        # TODO: handle a CSRF failure
        self.send_response(200)
        self.send_header("Content-type", "text/html")
        self.end_headers()

    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-type", "text/html")
        self.send_header("Access-Control-Allow-Origin", f"https://{BASE_DOMAIN}")
        # set cookies
        cookie = http.cookies.SimpleCookie()
        cookie['state'] = STATE
        cookie['url'] = f"https://{DOMAIN}/token"
        cookie['redir_url'] = REDIRECT_URL
        cookie['client_id'] = CLIENT_ID
        cookie['verifier'] = VERIFIER
        for morsel in cookie.values():
            self.send_header("Set-Cookie", morsel.OutputString())
        self.end_headers()
        # read a file and return it
        with open('redir.html','r') as htmlfile:
            html = htmlfile.read()
        self.wfile.write(html.encode('utf-8'))

def check_token(token):
    pass


if __name__ == "__main__":
    #
    wgc = Wg2faClient(REDIRECT_URL, BASE_DOMAIN, CLIENT_ID)
    # start out callback handler
    webServer = HTTPServer(('localhost', 8080), HttpCallbackHandler)
    # build challenge and login url
    VERIFIER = secrets.token_urlsafe(nbytes=32)
    challenge = Wg2faClient.get_challenge(VERIFIER)
    STATE = secrets.token_urlsafe(nbytes=16)
    url = wgc.get_login_url(challenge, STATE)
    print(url)
    # start the GUI
    # app = QApplication(sys.argv)
    # browser = Browser(url)
    # sys.exit(app.exec_())

    # open the page in a users web browser
    webbrowser.open_new_tab(url)
    # # start the callback handler
    try:
        webServer.serve_forever()
    except KeyboardInterrupt:
        pass