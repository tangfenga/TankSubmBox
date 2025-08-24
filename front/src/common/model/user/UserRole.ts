import ColorSelectionOption from '../base/option/ColorSelectionOption';
import Color from '../base/option/Color';
import Lang from '../global/Lang';

enum UserRole {
  GUEST = 'GUEST',
  USER = 'USER',
  ADMINISTRATOR = 'ADMINISTRATOR',
}

interface UserGroup {
  id: number;
  name: string;
  display: string;
  editable_labels: string;
  editable: boolean;
}

interface Label {
  id: number;
  name: string;
  type: 'bool' | 'numb';
}

let UserRoles: UserRole[] = Object.keys(UserRole).map((k) => k as UserRole);

let UserRoleMap: { [key in keyof typeof UserRole]: ColorSelectionOption } = {
  GUEST: {
    name: Lang.t('user.roleGuest'),
    value: 'GUEST',
    color: Color.WARNING,
  },
  USER: {
    name: Lang.t('user.roleUser'),
    value: 'USER',
    color: Color.PRIMARY,
  },
  ADMINISTRATOR: {
    name: Lang.t('user.roleAdministrator'),
    value: 'ADMINISTRATOR',
    color: Color.DANGER,
  },
};

let UserRoleList: ColorSelectionOption[] = [];
UserRoles.forEach((type: UserRole, index: number) => {
  UserRoleList.push(UserRoleMap[type]);
});

export { UserRole, UserRoles, UserRoleMap, UserRoleList };  export type { UserGroup, Label };

