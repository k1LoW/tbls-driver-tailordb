# ContactListMember

## Description

ContactListMember model

## Columns

| Name | Type | Default | Nullable | Children | Parents | Comment |
| ---- | ---- | ------- | -------- | -------- | ------- | ------- |
| id | uuid |  | false |  |  |  |
| contactId | uuid |  | false |  | [Contact](Contact.md) | Contact ID |
| contact | Contact |  | true |  | [Contact](Contact.md) | Link to the Contact |
| contactListId | uuid |  | false |  | [ContactList](ContactList.md) | ContactList ID |
| contactList | ContactList |  | true |  | [ContactList](ContactList.md) | Link to the ContactList |
| addedAt | string |  | true |  |  | Time when the contact was added to the list |
| createdAt | datetime |  | true |  |  | createdAt |
| updatedAt | datetime |  | true |  |  | updatedAt |

## Constraints

| Name | Type | Definition |
| ---- | ---- | ---------- |
| ForeignKey for contactList to ContactList | FOREIGN KEY | ForeignKeyType: ContactList |
| contactListCompositeKey | UNIQUE | {"contactListCompositeKey":{"FieldNames":["contactId","contactListId"],"Unique":true}} |

## Indexes

| Name | Definition |
| ---- | ---------- |
| Index for createdAt | Index: true |
| Index for updatedAt | Index: true |
| contactListCompositeKey | {"contactListCompositeKey":{"FieldNames":["contactId","contactListId"],"Unique":true}} |

## Relations

![er](ContactListMember.svg)

---

> Generated by [tbls](https://github.com/k1LoW/tbls)
