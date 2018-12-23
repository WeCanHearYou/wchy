import "./Avatar.scss";

import React from "react";
import { classSet, Fider } from "@fider/services";
import { isCollaborator, UserRole } from "@fider/models";

interface AvatarProps {
  user?: {
    id?: number;
    name: string;
    role?: UserRole;
    avatarURL?: string;
  };
  size?: "small" | "normal" | "large";
}

export const Avatar = (props: AvatarProps) => {
  const size = props.size || "normal";
  const id = props.user ? props.user.id : 0;
  const name = props.user ? props.user.name : "";
  const url = `${Fider.settings.tenantAssetsURL}/avatars/letter/${id}/${encodeURIComponent(name || "?")}`;
  const avatarURL = props.user ? props.user.avatarURL || url : url;

  const className = classSet({
    "c-avatar": true,
    [`m-${size}`]: true,
    "m-staff": props.user && props.user.role && isCollaborator(props.user.role)
  });

  return <img className={className} title={name} src={`${avatarURL}?size=50`} />;
};
