import { Star, Quote } from 'lucide-react'

const testimonials = [
  {
    id: 1,
    name: 'Sarah Johnson',
    role: 'Wedding Client',
    avatar: 'https://i.pravatar.cc/80?img=47',
    rating: 5,
    text: "The wedding cake was absolutely stunning — exactly what I envisioned. Our guests couldn't stop talking about how delicious it was. Cake Shop made our special day even more magical!",
  },
  {
    id: 2,
    name: 'Michael Chen',
    role: 'Birthday Customer',
    avatar: 'https://i.pravatar.cc/80?img=33',
    rating: 5,
    text: "I ordered a custom birthday cake for my daughter's 5th birthday. The castle design was perfect and the vanilla sponge was the fluffiest I've ever tasted. Will definitely order again!",
  },
  {
    id: 3,
    name: 'Emily Rodriguez',
    role: 'Regular Customer',
    avatar: 'https://i.pravatar.cc/80?img=5',
    rating: 5,
    text: "The ordering process was so easy and the delivery was right on time. The cheesecake was rich, creamy, and had the perfect crust. Cake Shop is my go-to for all celebrations now.",
  },
]

export function Testimonials() {
  return (
    <section className="py-20 bg-orange-50">
      <div className="container mx-auto max-w-6xl px-4">
        {/* Section header */}
        <div className="text-center mb-12">
          <p className="text-sm font-semibold text-primary uppercase tracking-widest mb-2">
            What Customers Say
          </p>
          <h2 className="text-4xl font-extrabold text-gray-900">Loved by Cake Lovers</h2>
          <p className="mt-3 text-muted-foreground max-w-md mx-auto">
            Join thousands of happy customers who celebrate life's sweetest moments with us.
          </p>
        </div>

        {/* Testimonials grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {testimonials.map((t) => (
            <div
              key={t.id}
              className="relative bg-white rounded-2xl p-6 shadow-sm hover:shadow-md transition-shadow"
            >
              <Quote className="absolute top-4 right-4 h-8 w-8 text-orange-100" />

              {/* Stars */}
              <div className="flex gap-1 mb-4">
                {Array.from({ length: t.rating }).map((_, i) => (
                  <Star key={i} className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                ))}
              </div>

              {/* Text */}
              <p className="text-sm text-gray-600 leading-relaxed mb-5 italic">"{t.text}"</p>

              {/* Author */}
              <div className="flex items-center gap-3">
                <img
                  src={t.avatar}
                  alt={t.name}
                  className="h-10 w-10 rounded-full object-cover"
                />
                <div>
                  <p className="font-semibold text-sm text-gray-900">{t.name}</p>
                  <p className="text-xs text-muted-foreground">{t.role}</p>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
