const Card = ({ items }) => {
  return (
    <article className="max-w-md mx-auto mt-4 shadow-lg border rounded-md duration-300 hover:shadow-sm">
      <a href={items.href} target="_blank">
        <img
          src={items.img}
          loading="lazy"
          alt={items.title}
          className="w-full h-48 rounded-t-md"
        />
        <div className="pt-3 ml-4 mr-2 mb-3">
          <h3 className="text-xl text-gray-900">{items.title}</h3>
          <p className="text-gray-400 text-sm mt-1">{items.desc}</p>
        </div>
      </a>
    </article>
  );
};

export default Card;
