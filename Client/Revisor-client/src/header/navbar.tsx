import { useGlobalContext } from "../app/context/GlobalContext";
import GoogleAuth from "../app/ui/auth/googleAuth";
import Logout from "../app/ui/auth/logout";
import "./navbar.css";
interface NavbarProps {
  logoUrl: string;
}
const NavBar:React.FC<NavbarProps> = ({logoUrl}) =>{
    const {email} = useGlobalContext();
    return (
        <>
        <nav className="navbar">
      <div className="navbar-left">
        <img src={logoUrl} alt="Logo" className="logo" />
        <span className="brand">Revisor</span>
      </div>
    </nav>
    {/* Signup or Signin button*/}
    {email ? <Logout /> :  <GoogleAuth/>}
        </>
    )
}
export default NavBar