import "./navbar.css";
interface NavbarProps {
  logoUrl: string;
}
const NavBar:React.FC<NavbarProps> = ({logoUrl}) =>{
    return (
        <>
        <nav className="navbar">
      <div className="navbar-left">
        <img src={logoUrl} alt="Logo" className="logo" />
        <span className="brand">Revisor</span>
      </div>
    </nav>
        </>
    )
}
export default NavBar