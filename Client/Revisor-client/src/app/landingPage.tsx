import NavBar from "../header/navbar";
import { useGlobalContext } from "./context/GlobalContext";
import Footer from "./footer/footer";
import Main from "./main/main";
import { useEffect } from "react";

const LandingPage = () =>{
  const {setEmail} = useGlobalContext(); //set value to global email state
    //Run hook to fetch code from query parameters
   useEffect(() => {
  const urlParams = new URLSearchParams(window.location.search);
  const code = urlParams.get('code');

  if (code) {
    fetch('http://localhost:8080/auth/google', {
      method: 'POST',
      credentials: "include",
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ "code":code }),
    })
      .then(res => res.json())
      .then(data => {
        //console.log('User data:', data.user);
        localStorage.setItem("name",data.user.name);
        localStorage.setItem("email",data.user.email);
        localStorage.setItem("token",data.token);
        localStorage.setItem("tokenExpiry",data.tokenExpiresAt);
        setEmail(data.user.email);
        console.log("User login successful");
      })
      .catch((error)=>{
        console.error("Failed to make post request to auth/google : ",error);
      })
    }
        //Now a timer will be started
        //This will trigger auto logout when the token expires
        const tokenExpiresAt = new Date((localStorage.getItem("tokenExpiry") as string)).getTime();
        const now = Date.now();
        const expiresInMs = tokenExpiresAt - now;
        //console.log(expiresInMs);
        if (expiresInMs > 0){
          setTimeout(()=>{
              //remove data from local storage
               localStorage.clear();
               setEmail('');
               console.log("Token expired and you are marked as Logged out");
          },expiresInMs);
        }
}, []);

    return (
        <>
        {/*NavBar header */}
        <NavBar logoUrl="broken"/>
        {/*Main section */}
        <Main/>
        {/*Footer section */}
        <Footer/>
        </>
    )
}
export default LandingPage