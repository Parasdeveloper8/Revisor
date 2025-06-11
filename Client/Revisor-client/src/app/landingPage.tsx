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
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ "code":code }),
    })
      .then(res => res.json())
      .then(data => {
        //console.log('User data:', data.user);
        localStorage.setItem("name",data.user.name);
        localStorage.setItem("email",data.user.email);
        localStorage.setItem("token",data.token.access_token);
        setEmail(data.user.email);
        console.log("User login successful");
      })
      .catch((error)=>{
        console.error("Failed to make post request to auth/google : ",error);
      })
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