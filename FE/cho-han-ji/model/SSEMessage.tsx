export class Message {
  MessageType: string

  constructor(messageType: string) {
    this.MessageType = messageType
  }
}

export class TypedMessage<T> extends Message {
  Message: T

  constructor(messageType: string, message: T) {
    super(messageType)
    this.MessageType = messageType
    this.Message = message
  }
}
