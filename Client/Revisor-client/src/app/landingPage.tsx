import NavBar from "./header/navbar";
import { useGlobalContext } from "./context/GlobalContext";
import Footer from "./footer/footer";
import Main from "./main/main";
import { useEffect } from "react";

const LandingPage = () =>{
  const {setUserData,tokenExpiry,setInfo} = useGlobalContext(); //set value to global email state
    //Run hook to fetch code from query parameters
   useEffect(() => {
   const urlParams = new URLSearchParams(window.location.search);
   const code = urlParams.get('code');

           //send request to /auth/me
        fetch('http://localhost:8080/auth/me', {
         method: 'GET',
         credentials: "include",
        })
         .then(res => res.json())
         .then(data => {
        const userData = {
           name: data.user.name,
           email: data.user.email,
           token: data.token,
           tokenExpiry: new Date(data.tokenExpiresAt), // convert to Date
        };
        setUserData(userData);
        console.log("User's data fetched successfully");
        setInfo({message:"User login successful",anyInfo:"success"});
        setTimeout(()=>setInfo({message:'',anyInfo:'noInfo'}),2500);
        })
        .catch((error)=>{
           setUserData({name:'',email:'',token:'',tokenExpiry:new Date});
           console.log("Error occurred.Login process will start: ",error);
           //Now if there is any error it tells that there is something wrong with server or User is not logged in
           //Let's assume that user is not logged in 
           //Then call /auth/google api to do login process
            fetch('http://localhost:8080/auth/google', {
             method: 'POST',
             credentials: "include",
             headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({ "code":code }),
              })
               .then(res => res.json())
               .then(data => {
              //console.log('User data:', data.user);
        const userData = {
        name: data.user.name,
        email: data.user.email,
        token: data.token,
        tokenExpiry: new Date(data.tokenExpiresAt), // convert to Date
        };
                  setUserData(userData);
                   console.log("User login successful");
                   setInfo({message:"User login successful",anyInfo:"success"});
                   setTimeout(()=>setInfo({message:'',anyInfo:'noInfo'}),2500);
                  })
              .catch((error)=>{
                   console.error("Error occurred: ",error);
                   setInfo({message:"User login failed",anyInfo:"error"});
                   setTimeout(()=>setInfo({message:'',anyInfo:'noInfo'}),2500);
                 })
        });
  
          //function to handle logout
         const handleLogout = ()=>{
           setUserData({name:'',email:'',token:'',tokenExpiry:new Date});
           console.log("Token expired and you are marked as Logged out");
           setInfo({message:"User logout successful",anyInfo:"success"});
           setTimeout(()=>setInfo({message:'',anyInfo:'noInfo'}),2500);
         }
        //Now a timer will be started
        //This will trigger auto logout when the token expires
        const tokenExpiresAt = tokenExpiry.getTime();
        const now = Date.now();
        const expiresInMs = tokenExpiresAt - now;
        //console.log(expiresInMs);
        if (expiresInMs > 0){
          setTimeout(()=>{
              handleLogout();
          },expiresInMs);
        }
        //If current time exceeds token expiry
        if (tokenExpiresAt <= now){
            handleLogout();
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