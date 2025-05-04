"use client"
import React from "react";
import {
  Collapse,
  Container,
  Nav,
  Navbar,
  NavbarToggler,
  NavItem,
} from "reactstrap";
import {NavLink} from "react-router-dom";
import Profile from "./Profile";

import UserInfo from "@/lib/UserInfo";
import {SSOEnabled} from "@/lib/Configuration";
import {
  getUseOidcAccessToken,
  NoSSOUserInfo,
} from "@/lib/SSODisabledDefaults";

const NavBar: React.FunctionComponent = () => {
  const [isOpen, setIsOpen] = React.useState<boolean>(false);

  const toggle = () => {
    setIsOpen(!isOpen);
  };
  const {accessTokenPayload} = getUseOidcAccessToken()();
  const userInfo = SSOEnabled
    ? (accessTokenPayload as UserInfo)
    : NoSSOUserInfo;
  
  return (
    <div>
      <Navbar color="primary" dark expand="lg">
        <Container>
          <NavLink to="/" className={"navbar-brand"}>
            InfoSys
          </NavLink>
          <NavbarToggler onClick={toggle} />
          <Collapse isOpen={isOpen} navbar>
            <Nav navbar>
              {/* to add stuff to the navbar, add a NavItem tag with a NavLink to the route */}
              <NavItem>
                <NavLink to="/" className="nav-link">
                  Home
                </NavLink>
              </NavItem>
              <NavItem>
                <NavLink to="/message" className="nav-link">
                  Submit a Message
                </NavLink>
              </NavItem>
              <NavItem>
                <NavLink to={`message/${userInfo.preferred_username}`} className="nav-link">
                  My Messages
                </NavLink>
              </NavItem>
            </Nav>
            <Nav navbar className="ml-auto">
              <Profile />
              {/*
                FIXME: ThemeToggle from InfoSys-Frontend depends on Helmet, which is giving me
                conflicts.
                https://github.com/ikopke23/InfoSys-Frontend/blob/master/src/components/ThemeToggle.tsx 
              <NavItem className="nav-link">
                <ThemeToggle />
              </NavItem>*/}
            </Nav>
          </Collapse>
        </Container>
      </Navbar>
    </div>
  );
};

export default NavBar;
