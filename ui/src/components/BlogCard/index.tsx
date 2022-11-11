const Card = ({ items }) => {
  return (
    <a href={items.blog_url} target="_blank">
      {items.blog_name}
    </a>
  );
};

export default Card;
