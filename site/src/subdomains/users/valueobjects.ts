import { VO } from './../../common/vo/vo';

// name
// password
// login

const nameRegex = new RegExp("^[a-zA-Z0-9-_]*$")

const toShortNameError = new Error("name has to be at least 3 characters long")
const toLongNameError = new Error("name cannot have more than 64 characters")
const invalidCharachtersNameError = new Error("name can contain only letters, numbers, - and _")

export class Name implements VO {
  public Value: string

  constructor(value: string) { this.Value = value.toLowerCase().trim() }

  Valid(): Error[] {
    let errs: Error[] = []
    if (this.Value.length < 3) { errs.push(toShortNameError) }
    if (this.Value.length > 64) { errs.push(toLongNameError) }
    if (!nameRegex.test(this.Value)) { errs.push(invalidCharachtersNameError) }
    return errs
  }

  Errors(): Error[] {
    return [toShortNameError, toLongNameError, invalidCharachtersNameError]
  }
}

// password

const smallLeterRegex = new RegExp("[a-z]")
// const bigLeterRegex = new RegExp("[A-Z]")
const numberRegex = new RegExp("[0-9]")
// const specialCharacterRegex = new RegExp("[!@#$%^&*()_+{}:\"|<>?]")


const toShortPasswordError = new Error("password has to be at least 3 characters long")
const toLongPasswordError = new Error("password cannot have more than 64 characters")
const noSmallCharactersPasswordError = new Error("password must contain small letter")
// const noBigCharactersPasswordError = new Error("password must contain big letter")
const noNumbersPasswordError = new Error("password must contain number")
// const noSpecialCharactersPasswordError = new Error("password must contain special character")

// export const PasswordInputProps: { [key: string]: any } = {
//   type: "password",
//   placeholder: "********",
// }
export class Password implements VO {
  public Value: string

  constructor(value: string) {
    this.Value = value
  }

  Valid(): Error[] {
    let errs: Error[] = []
    if (this.Value.length < 3) { errs.push(toShortPasswordError) }
    if (this.Value.length > 64) { errs.push(toLongPasswordError) }
    if (!smallLeterRegex.test(this.Value)) { errs.push(noSmallCharactersPasswordError) }
    // if (!bigLeterRegex.test(this.Value)) { errs.push(noBigCharactersPasswordError) }
    if (!numberRegex.test(this.Value)) { errs.push(noNumbersPasswordError) }
    // if (!specialCharacterRegex.test(this.Value)) { errs.push(noSpecialCharactersPasswordError) }
    return errs
  }

  Errors(): Error[] {
    return [
      toShortPasswordError,
      toLongPasswordError,
      noSmallCharactersPasswordError,
      // noBigCharactersPasswordError,
      noNumbersPasswordError,
      // noSpecialCharactersPasswordError,
    ]
  }
}

// login

export class Login implements VO {
  public Value: string

  constructor(value: string) {
    this.Value = new Name(value).Value
  }

  Valid(): Error[] {
    return new Name(this.Value).Valid()
  }

  Errors(): Error[] {
    return new Name(this.Value).Errors()
  }
}