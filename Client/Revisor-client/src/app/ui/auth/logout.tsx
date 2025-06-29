
import { useGlobalContext } from "../../context/GlobalContext";
import "./logout.css";

const Logout = ()=>{
  // console.log(localStorage.getItem("token"));
    const {setUserData,token} = useGlobalContext(); //Set empty value to global email state

    //function to logout user
    const handleLogout = () =>{
         fetch('http://localhost:8080/auth/logout',{
          method: 'POST',
          credentials: "include",
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({"token":token}),
         })
         .then(response =>{
            if(response.ok){
                setUserData({name:'',email:'',token:'',tokenExpiry:new Date});
               console.log("Logged out");
            }
         })
         .catch((error) => {
            console.error("Failed to make post request to auth/logout: ",error);
         })
    }
    return (
        <>
        <button onClick={handleLogout} className="logout-btn">Logout</button>
        </>
    )
}
export default Logout