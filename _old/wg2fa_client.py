from hashlib import sha256
from base64 import urlsafe_b64encode
from urllib.parse import urlparse, parse_qs, quote_plus

class Wg2faClient:
    def __init__(self, redirect_url, base_domain, client_id):
        self.redirect_url = redirect_url
        self.base_domain = base_domain
        self.client_id = client_id
        self.domain = f'{self.base_domain}/oauth2/default/v1'

    @staticmethod
    def get_challenge(verifier):
        hasher = sha256()
        hasher.update(verifier.encode('utf-8'))
        encoded = urlsafe_b64encode(hasher.digest()).decode('utf-8')
        encoded = encoded.replace('=','')
        return encoded
    
    @staticmethod
    def get_

    def get_login_url(self, challenge, state):
        eurl = quote_plus(self.redirect_url)
        url = f"https://{self.domain}/authorize?"
        url = url + "response_type=code"
        url = url + f"&client_id={self.client_id}"
        url = url + f"&redirect_uri={eurl}"
        url = url + f"&response_mode=fragment"
        url = url + f"&code_challenge={challenge}"
        url = url + f"&code_challenge_method=S256" 
        #url = url + "&scope=openid%20profile%20email&audience=appointments:api&state=xyzABC123"
        url = url + f"&scope=openid%20profile%20email&state={state}"
        return url