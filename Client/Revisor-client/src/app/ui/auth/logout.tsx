
import { useGlobalContext } from "../../context/GlobalContext";
import "./logout.css";

const Logout = ()=>{
  // console.log(localStorage.getItem("token"));
    const {setEmail} = useGlobalContext(); //Set empty value to global email state

    //function to logout user
    const handleLogout = () =>{
         const token = localStorage.getItem("token");
         fetch('http://localhost:8080/auth/logout',{
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({"token":token}),
         })
         .then(response =>{
            if(response.ok){
               //remove data from local storage
               localStorage.removeItem("name");
               localStorage.removeItem("email");
               localStorage.removeItem("token");
               setEmail('');
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