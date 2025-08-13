import { useGlobalContext } from "../context/GlobalContext";
import GoogleAuth from "../ui/auth/googleAuth";
import Logout from "../ui/auth/logout";
import NotificationBar from "../ui/notificationBar/notificationbar";
import "./navbar.css";
interface NavbarProps {
  logoUrl: string;
}
const NavBar:React.FC<NavbarProps> = ({logoUrl}) =>{
    const {email,info} = useGlobalContext();
    return (
        <>
        <nav className="navbar">
      <div className="navbar-left">
        <img src={logoUrl} alt="Logo" className="logo" />
        <span className="brand">Revisor</span>
      </div>
    </nav>
    {/*Notification Bar*/}
    {
      info.anyInfo === "success" ? (
         <NotificationBar message={info.message} type="success" duration={2500}/>
      ) :
      info.anyInfo === "error" ? (
         <NotificationBar message={info.message} type="error" duration={2500}/>
      ) : null
    }
    {/* Signup or Signin button*/}
    {email ? <Logout /> :  <GoogleAuth/>}
        </>
    )
}
export default NavBar