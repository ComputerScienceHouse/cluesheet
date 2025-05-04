import {
  DropdownItem,
  DropdownMenu,
  DropdownToggle,
  UncontrolledDropdown,
} from "reactstrap";

import React from "react";
import UserInfo from "@/lib/UserInfo";
import {SSOEnabled} from "@/lib/Configuration";
import {
  getUseOidcAccessToken,
  getUseOidcHook,
  NoSSOProfilePicture,
  NoSSOUserInfo,
} from "@/lib/SSODisabledDefaults";

const Profile: React.FunctionComponent = () => {
  const {logout} = getUseOidcHook()();
  const {accessTokenPayload} = getUseOidcAccessToken()();
  const userInfo = SSOEnabled
    ? (accessTokenPayload as UserInfo)
    : NoSSOUserInfo;

  return (
    <UncontrolledDropdown nav inNavbar>
      <DropdownToggle nav caret className="navbar-user">
        <img
          className="rounded-circle"
          src={
            SSOEnabled
              ? `https://profiles.csh.rit.edu/image/${userInfo.preferred_username}`
              : NoSSOProfilePicture
          }
          alt=""
          aria-hidden={true}
          width={32}
          height={32}
        />
        ({userInfo.preferred_username})
        <span className="caret" />
      </DropdownToggle>
      <DropdownMenu>
        {
          // to add stuff to the profile dropdown, you can
        }
        <DropdownItem href="https://members.csh.rit.edu">Members</DropdownItem>
        <DropdownItem divider />
        <DropdownItem onClick={() => logout(null)}>Logout</DropdownItem>
      </DropdownMenu>
    </UncontrolledDropdown>
  );
};

export default Profile;
