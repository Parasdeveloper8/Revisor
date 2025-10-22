// type to hold each flashcard item
export type FlashCardItem = {
  Email: string;
  TopicName: string;
  Time: string;
  FormattedTime: string;
  Data: {
    Heading: string;
    Value: string;
  }[];
  Uid : string;
};
