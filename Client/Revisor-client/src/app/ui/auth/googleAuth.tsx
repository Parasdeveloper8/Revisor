import "./googleAuth.css";
const CLIENT_ID = "596784384543-v8mcuvrtklis1hhp1e3ici4rkoemhf7i.apps.googleusercontent.com";
const REDIRECT_URI = 'http://localhost:5173';

const GoogleAuth = () =>{
  //function to handle redirection to consent screen
  const handleLogin = () => {
    const googleOAuthURL = `https://accounts.google.com/o/oauth2/v2/auth?client_id=${CLIENT_ID}&redirect_uri=${REDIRECT_URI}&response_type=code&scope=openid%20email%20profile&access_type=offline&prompt=consent`;
    window.location.href = googleOAuthURL;
  };
    return (
    <button onClick={handleLogin} className="sign-btn">
      <img
        src="https://developers.google.com/identity/images/g-logo.png"
        alt="Google" className="google-image"
      />
      Sign in with Google
    </button>
  );
}
export default GoogleAuth